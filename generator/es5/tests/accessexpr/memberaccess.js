$module('memberaccess', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(2, $g.____testlib.basictypes.Integer)).then(function (result) {
        instance.someInt = result;
      }));
      init.push($promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.someBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Build = function () {
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $g.memberaccess.SomeClass.new().then(function ($result0) {
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
    $instance.InstanceFunc = function () {
      var $this = this;
      return $promise.empty();
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this.someInt);
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Build', 1, $g.____testlib.basictypes.Function($g.memberaccess.SomeClass).$typeref()], ['InstanceFunc', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['SomeProp', 3, $g.____testlib.basictypes.Integer.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.memberaccess.SomeClass).$typeref()]);
    };
  });

  $static.DoSomething = function (sc, scn) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.dynamicaccess($g.memberaccess.SomeClass, 'Build');
            $t.dynamicaccess($g.memberaccess.SomeClass, 'Build');
            $t.dynamicaccess(maimport, 'AnotherFunction');
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
            sc.InstanceFunc().then(function ($result0) {
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
            $t.dynamicaccess(sc, 'InstanceFunc');
            sc.SomeProp().then(function ($result0) {
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
            sc.SomeProp().then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            scn.SomeProp().then(function ($result0) {
              $result = $result0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
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
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.memberaccess.SomeClass.new().then(function ($result0) {
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
            $promise.resolve($t.unbox(sc.someBool)).then(function ($result0) {
              $result = $t.box($result0 && $t.unbox(sc.someBool), $g.____testlib.basictypes.Boolean);
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
