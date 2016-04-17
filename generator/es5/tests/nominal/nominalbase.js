$module('nominalbase', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.SomeField = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  this.$type('FirstNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(!$t.unbox($t.unbox($this).SomeField), $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
  });

  this.$type('SecondNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $t.box($this, $g.nominalbase.FirstNominal).SomeProp().then(function ($result0) {
                $result = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
                $current = 1;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 1:
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
  });

  $static.TEST = function () {
    var sc;
    var sn;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.nominalbase.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            sn = $t.box($t.box(sc, $g.nominalbase.FirstNominal), $g.nominalbase.SecondNominal);
            sn.GetValue().then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
