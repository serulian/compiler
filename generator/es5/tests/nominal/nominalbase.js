$module('nominalbase', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
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
    this.new = function ($wrapped) {
      var instance = new this();
      instance.$wrapped = $wrapped;
      return instance;
    };
    this.$apply = function (data) {
      var instance = new this();
      instance.$wrapped = data.$wrapped;
      return instance;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.nominalwrap(!$t.nominalunwrap($t.nominalunwrap($this).SomeField), $g.____testlib.basictypes.Boolean));
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    });
  });

  this.$type('SecondNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.new = function ($wrapped) {
      var instance = new this();
      instance.$wrapped = $wrapped;
      return instance;
    };
    this.$apply = function (data) {
      var instance = new this();
      instance.$wrapped = $g.nominalbase.FirstNominal.$apply(data.$wrapped);
      return instance;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.nominalunwrap($this).SomeProp().then(function ($result0) {
                $result = $t.nominalwrap(!$t.nominalunwrap($result0), $g.____testlib.basictypes.Boolean);
                $state.current = 1;
                $callback($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var sc;
    var sn;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.nominalbase.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            sn = $t.nominalwrap($t.nominalwrap(sc, $g.nominalbase.FirstNominal), $g.nominalbase.SecondNominal);
            sn.GetValue().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
