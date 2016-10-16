$module('shortcut', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Value', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.shortcut.SomeClass).$typeref()]);
    };
  });

  this.$type('SomeNominalType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.shortcut.SomeClass;
    };
    this.$typesig = function () {
      return $t.createtypesig();
    };
  });

  $static.DoSomething = function (sc) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            sc.Value().then(function ($result0) {
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
  $static.TEST = function () {
    var $result;
    var st;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.shortcut.SomeClass.new().then(function ($result0) {
              $result = $t.box($result0, $g.shortcut.SomeNominalType);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            st = $result;
            $g.shortcut.DoSomething($t.unbox(st)).then(function ($result0) {
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
